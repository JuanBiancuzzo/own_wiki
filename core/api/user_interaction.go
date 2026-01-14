package api

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"

	db "github.com/JuanBiancuzzo/own_wiki/core/database"
	c "github.com/JuanBiancuzzo/own_wiki/core/systems/configuration"

	pb "github.com/JuanBiancuzzo/own_wiki/core/api/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type userInteraction struct {
	pb.UnimplementedUserInteractionServer
}

func newUserInteraction() *userInteraction {
	return &userInteraction{}
}

func (*userInteraction) LoadPlugin(ctx context.Context, loadPluginRequest *pb.LoadPluginRequest) (*pb.LoadPluginResponse, error) {
	// We would need to find the files (build it as a plugin), then load it and keep the reference to that plugin.
	// With the plugin, we should call the necesary function to define the structures given by the user. This should
	// be the components, and the views. Then the components should be return

	return nil, nil
}

func (*userInteraction) ImportFiles(importFileStream pb.UserInteraction_ImportFilesServer) error {
	// We could send via a channel the file paths needed, and via a gorouting be processing them as the
	// user defines. Finally, via another channel, send the component description to the stream

	return nil
}

func (*userInteraction) Render(renderStream pb.UserInteraction_RenderServer) error {
	// First, we would need to create the SystemInteractionClient to get the Query and SendEvent functions.
	// We would need the view register in the LoadPlugin function, and create the gorouting to continuously
	// be able to send events, and get the scene render

	return nil
}

func (*userInteraction) Close() {
	// Close the plugin
}

type UserInteractionServer struct {
	listener net.Listener
	server   *grpc.Server
}

func NewUserInteractionServer(config c.UserInteractionConfig) (*UserInteractionServer, error) {
	direction := fmt.Sprintf("%s:%d", config.Ip, config.Port)
	lis, err := net.Listen(config.Protocol, direction)
	if err != nil {
		return nil, fmt.Errorf("failed to create listener server for UserInteraction at %s, with error: %v", direction, err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserInteractionServer(grpcServer, newUserInteraction())

	return &UserInteractionServer{
		listener: lis,
		server:   grpcServer,
	}, nil
}

func (us *UserInteractionServer) Serve() error {
	if err := us.server.Serve(us.listener); err != nil {
		return fmt.Errorf("failed to serve UserInteraction, with error: %v", err)
	}
	return nil
}

type UserInteractionClient struct {
	Conn *grpc.ClientConn
	User pb.UserInteractionClient
}

func NewUserInteractionClient(config c.UserInteractionConfig) (*UserInteractionClient, error) {
	transportCredentials := grpc.WithTransportCredentials(insecure.NewCredentials())
	direction := fmt.Sprintf("%s:%d", config.Ip, config.Port)

	conn, err := grpc.NewClient(direction, transportCredentials)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client for UserInteraction at %s, with error: %v", direction, err)
	}

	return &UserInteractionClient{
		Conn: conn,
		User: pb.NewUserInteractionClient(conn),
	}, nil
}

type seenType int

const (
	ST_NOT_SEEN = iota
	ST_SEEN
	ST_COMPLETED
)

type seenComponent struct {
	Component *pb.ComponentStructure
	Seen      seenType
}

func dfs(name string, tableNames map[string]*seenComponent, exec func(*pb.ComponentStructure) error) error {
	seenComponent := tableNames[name]
	seenComponent.Seen = ST_SEEN

	component := seenComponent.Component
	for _, field := range component.GetFields() {
		if _, ok := field.GetType().(*pb.FieldStructure_Reference); !ok {
			continue
		}
		refTableName := field.GetReference().GetTableName()
		refComponent := tableNames[refTableName]

		switch refComponent.Seen {
		case ST_NOT_SEEN:
			if err := dfs(refTableName, tableNames, exec); err != nil {
				return err
			}

		case ST_SEEN:
			return fmt.Errorf("In the table %s got a cycle", name)
		}
	}

	seenComponent.Seen = ST_COMPLETED
	return exec(component)
}

func (uc *UserInteractionClient) LoadPlugin(ctx context.Context, pluginPath string) (descriptions []db.TableDescription, err error) {
	response, err := uc.User.LoadPlugin(ctx, &pb.LoadPluginRequest{PluginPath: pluginPath})
	if err != nil {
		return descriptions, fmt.Errorf("Failed to load plugin at %q and get component data, with error: %v", pluginPath, err)
	}

	tableNames := make(map[string]*seenComponent)
	for _, component := range response.GetComponents() {
		name := component.GetName()
		if _, ok := tableNames[name]; ok {
			return descriptions, fmt.Errorf("Failed to create tables, there are multiple tables with the same name (atleast %s)", name)
		}
		tableNames[name] = &seenComponent{Component: component, Seen: ST_NOT_SEEN}
	}

	nameTables := make(map[string]*db.TableDescription)
	exec := func(component *pb.ComponentStructure) error {
		name := component.GetName()
		fieldStructures := component.GetFields()
		fields := make([]db.Field, len(fieldStructures))

		for i, fieldStructure := range fieldStructures {
			fieldName := fieldStructure.GetName()
			fieldIsNull := fieldStructure.GetIsNull()
			fieldIsKey := fieldStructure.GetIsKey()

			switch value := fieldStructure.GetType().(type) {
			case *pb.FieldStructure_Primitive:
				fieldType, err := value.Primitive.GetDataBaseFieldType()
				if err != nil {
					return fmt.Errorf("Primitive type failed, with error: %v", err)
				}
				fields[i] = db.NewPrimitiveField(fieldName, fieldType, fieldIsNull, fieldIsKey)

			case *pb.FieldStructure_Reference:
				reference := nameTables[value.Reference.GetTableName()]
				fields[i] = db.NewReferencesField(fieldName, reference, fieldIsNull, fieldIsKey)
			}
		}

		nameTables[name] = db.NewTableDescription(name, fields)
		return nil
	}

	for _, component := range tableNames {
		if component.Seen == ST_COMPLETED {
			continue
		}

		if err := dfs(component.Component.GetName(), tableNames, exec); err != nil {
			return descriptions, fmt.Errorf("Failed to create tables, with error: %v", err)
		}
	}

	descriptions = make([]db.TableDescription, 0, len(nameTables))
	for _, table := range nameTables {
		descriptions[len(descriptions)] = *table
	}

	return descriptions, nil
}

// TODO: Change this type to be the correct representation of the ComponentDescription, this should be able to be save in the database
type EntityDescription *pb.ComponentDescription

func (uc *UserInteractionClient) ImportFiles(ctx context.Context, sendFilePaths chan string, receiveEntity chan EntityDescription) error {
	stream, err := uc.User.ImportFiles(ctx)
	if err != nil {
		// We close the channel because there is no entity to be send
		close(receiveEntity)

		// We consume all the files send
		for range sendFilePaths {
		}

		return fmt.Errorf("Failed to create ImportFiles stream, with error: %v", err)
	}

	var waitSendAndReceive sync.WaitGroup
	errorChannel := make(chan error, 2)

	waitSendAndReceive.Add(1)
	go func(receiveFiles chan string, stream pb.UserInteraction_ImportFilesClient, wg *sync.WaitGroup) {
		errorOccurred := false

		for filePath := range receiveFiles {
			if errorOccurred {
				// We need to consume all the file send
				continue
			}

			if err := stream.Send(&pb.ImportedFilesRequest{FilePath: filePath}); err != nil {
				errorOccurred = true
				errorChannel <- fmt.Errorf("Error while sending file, with error: %v", err)
			}
		}

		if !errorOccurred {
			errorChannel <- nil
		}

		stream.CloseSend()
		wg.Done()
	}(sendFilePaths, stream, &waitSendAndReceive)

	waitSendAndReceive.Add(1)
	go func(sendEntity chan EntityDescription, stream pb.UserInteraction_ImportFilesClient, wg *sync.WaitGroup) {
		for {
			if response, err := stream.Recv(); err == io.EOF {
				errorChannel <- nil
				break

			} else if err != nil {
				errorChannel <- fmt.Errorf("Error while receiving entity information, with error: %v", err)
				break

			} else {
				sendEntity <- response.GetComponent()
			}
		}

		close(sendEntity)
		wg.Done()
	}(receiveEntity, stream, &waitSendAndReceive)

	firstError, secondError := <-errorChannel, <-errorChannel
	waitSendAndReceive.Wait()

	if firstError != nil && secondError != nil {
		return fmt.Errorf("Got the errors: %v, and %v", firstError, secondError)

	} else if firstError != nil {
		return firstError

	} else if secondError != nil {
		return secondError
	}

	return nil
}

// TODO: Change the function to accept a channel for events (as in core/events/Event) to send, and a channel to get
// the scene representation
func (uc *UserInteractionClient) Render(ctx context.Context) (pb.UserInteraction_RenderClient, error) {
	return uc.User.Render(ctx)
}

func (uc *UserInteractionClient) Close() {
	uc.Conn.Close()
}

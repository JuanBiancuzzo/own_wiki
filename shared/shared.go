package shared

// Represents the structure of a component, this defines the structure of the data to be save
type ComponentStructure any

// Represents which views would be shown, and what information (entity) is need it to show it
type ViewInformation any

// Represents a component or a composition of components (an entity) that holds information
type EntityDescription any

// Is the string representation of a file, and its metadata
type File string

/*
This interface lets the user define the components, each entity and view for the project.
*/
type UserDefineStructure interface {
	// The components are the smallest data storage given by the system. They can depende on
	// each other, but there has to be a way to constructe them with out an infinite loop
	RegisterComponents() []ComponentStructure

	// Views are the representation of a entity to be shown by the program in the platform
	// define at compilation time
	RegisterViews() (mainViews ViewInformation, otherViews []ViewInformation)

	// Given that when importing file there has to be a way to transform them in entities, this
	// is where it happends. This also defines what entity is it wanted to be the main menu. If
	// multiples entities are main menu capable, then it will apear an option to select
	ProcessFile(file File) []EntityDescription
}

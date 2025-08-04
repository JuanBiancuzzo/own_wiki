FROM node:20-bookworm

# Instalando Tree-sitter CLI
RUN npm install -g tree-sitter-cli

# Copiando todo
COPY . .

# Haciendo ejecutable el bash script
RUN chmod +x /crear_bindings.sh

# Creando los bindings
CMD [ "./crear_bindings.sh" ]
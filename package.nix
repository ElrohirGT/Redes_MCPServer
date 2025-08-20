{
  lib,
  buildGoModule,
}:
buildGoModule {
  name = "Redes_MCPServer";
  src = ./.;
  vendorHash = null;
  meta = {
    description = "FAGD MCP Server for Redes course";
    homepage = "https://github.com/ElrohirGT/Redes_MCPServer";
    license = lib.licenses.mit;
    maintainers = with lib.maintainers; [elrohirgt];
  };
}

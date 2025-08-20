{
  lib,
  buildGoModule,
}:
buildGoModule {
  src = ./.;
  meta = {
    description = "FAGD MCP Server for Redes course";
    homepage = "https://github.com/ElrohirGT/Redes_MCPServer";
    license = lib.licenses.MIT;
    maintainers = with lib.maintainers; [elrohirgt];
  };
}

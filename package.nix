{
  lib,
  buildGoModule,
}:
buildGoModule {
  name = "Redes_MCPServer";
  src = ./.;
  vendorHash = "sha256-fTP/PZXcJUuDx3OA2zJSTqGTwcIAJI7qXeWlCit9f+k=";
  meta = {
    description = "FAGD MCP Server for Redes course";
    homepage = "https://github.com/ElrohirGT/Redes_MCPServer";
    license = lib.licenses.mit;
    maintainers = with lib.maintainers; [elrohirgt];
  };
}

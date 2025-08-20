{
  lib,
  buildGoModule,
}:
buildGoModule {
  name = "Redes_MCPServer";
  src = ./.;
  # vendorHash = "sha256-fTP/PZXcJUuDx3OA2zJSTqGTwcIAJI7qXeWlCit9f+k=";
  vendorHash = "sha256-8Vs08PkqO/voV08HThCr7lwQ7AeEa501qhgBTwWlxB4=";
  doCheck = false;
  meta = {
    description = "FAGD MCP Server for Redes course";
    homepage = "https://github.com/ElrohirGT/Redes_MCPServer";
    license = lib.licenses.mit;
    maintainers = with lib.maintainers; [elrohirgt];
  };
}

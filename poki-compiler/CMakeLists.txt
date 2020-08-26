set(GO_BINARY pokic)

add_custom_command(
  OUTPUT
    ${GO_BINARY}
  COMMAND
    go build -o ${CMAKE_CURRENT_BINARY_DIR}/${GO_BINARY}
  WORKING_DIRECTORY
    ${CMAKE_CURRENT_SOURCE_DIR}
)

add_custom_target(PokiCompiler ALL DEPENDS ${GO_BINARY})

install(
  FILES
    ${CMAKE_CURRENT_BINARY_DIR}/pokic
  DESTINATION
    ${CMAKE_INSTALL_BINDIR}
)

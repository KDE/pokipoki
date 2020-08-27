include_directories(headers)

add_library(pokipoki
  SHARED
    Database.cpp
)

install(
  TARGETS
    pokipoki
  DESTINATION
    ${CMAKE_INSTALL_LIBDIR}
)

target_link_libraries(
  pokipoki
    Qt5::Core
    Qt5::Sql
)

configure_file("pokipoki.pc.in" "pokipoki.pc" @ONLY)

install(
  FILES
    ${CMAKE_CURRENT_BINARY_DIR}/pokipoki.pc
  DESTINATION
    ${CMAKE_INSTALL_LIBDIR}/pkgconfig
)


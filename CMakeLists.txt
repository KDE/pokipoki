cmake_minimum_required(VERSION 3.5)

project(pokipoki LANGUAGES CXX)

set(CMAKE_INCLUDE_CURRENT_DIR ON)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

include(GNUInstallDirs)

find_package(Qt5 COMPONENTS Core Sql REQUIRED)

add_subdirectory(poki-compiler)
add_subdirectory(libpokipoki)
# add_subdirectory(example)

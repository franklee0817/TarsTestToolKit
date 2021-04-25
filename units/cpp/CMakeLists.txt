cmake_minimum_required(VERSION 2.8)


cmake_minimum_required(VERSION 2.8)

project(Demo-CppTars)

option(TARS_MYSQL "option for mysql" ON)
option(TARS_SSL "option for ssl" OFF)
option(TARS_GMSSL "option for ssl" OFF)
option(TARS_HTTP2 "option for http2" OFF)

if(WIN32)
include (c:\\tars\\cpp\\makefile\\tars-tools.cmake)
else()
include (/usr/local/tars/cpp/makefile/tars-tools.cmake)
endif()


set(CMAKE_BUILD_TYPE "Debug")
# set(CMAKE_BUILD_TYPE "Release")

set(TARS_WEB_HOST "" CACHE STRING "set web host")
IF (TARS_WEB_HOST STREQUAL "")
	set(TARS_WEB_HOST "http://taf.test.com")
ENDIF ()

include_directories(/usr/local/tars/cpp/thirdparty/include)
link_directories(/usr/local/tars/cpp/thirdparty/lib)

#include_directories(/home/tafjce/App/OtherServer)

add_subdirectory(src)

# target_link_libraries(TestServer ${LIB_MYSQL}) 


# gen_server(Demo CppTars)

#target_link_libraries(mysqlclient ssl crypto nghttp2-static)

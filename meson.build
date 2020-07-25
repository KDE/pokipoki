project('pokipoki', 'cpp', default_options: ['cpp_std=c++17'])

qt5 = import('qt5')
qt5_deps = dependency('qt5', modules: ['Core', 'Sql'])

pokipoki_headers = include_directories('libpokipoki/headers')

subdir('poki-compiler')
subdir('libpokipoki')
subdir('example')
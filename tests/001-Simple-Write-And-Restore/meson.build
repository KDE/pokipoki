moc_files = qt5.preprocess(
  moc_headers: '001.h',
  include_directories: pokipoki_headers,
)

eA = executable(
    '001-A',
    '001-A.cpp',
    moc_files,
    link_with: pokipoki_lib,
    dependencies: qt5_deps,
    include_directories: pokipoki_headers,
)

eB = executable(
    '001-B',
    '001-B.cpp',
    moc_files,
    link_with: pokipoki_lib,
    dependencies: qt5_deps,
    include_directories: pokipoki_headers,
)

test('001: Simple Write & Restore: Write', eA)
test('001: Simple Write & Restore: Restore', eB)

poki_gen = generator(
    find_program('pokic'),
    output: '@BASENAME@.c',
    arguments: ['-input', '@INPUT@', '-output', '@OUTPUT@']
)

tests = [
    '001-Simple-Write-And-Restore',
]

foreach test : tests
    subdir(test)
endforeach

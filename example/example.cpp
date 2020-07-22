#include <QCoreApplication>
#include "example.h"

int main(int argc, char *argv[])
{
    auto app = new QCoreApplication(argc, argv);
    app->setApplicationName("fimbeb-example");

    auto note = Note{};
    note.id = "pink";
    note.text = "yeet";

    FBDatabase::instance()->save(note);

    return 0;
}

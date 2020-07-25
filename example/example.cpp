#include <QCoreApplication>
#include "example.gen.h"

int main(int argc, char *argv[])
{
    auto app = new QCoreApplication(argc, argv);
    app->setApplicationName("fimbeb-example");

    auto note = Note::newNote();
    note->set_title("yeet");
    note->commit();

    auto child = Note::newNote();
    child->set_title("ohno");
    child->commit();

    note->addChildNote(child);

    return 0;
}

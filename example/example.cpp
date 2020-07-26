#include <QCoreApplication>
#include "example.gen.h"

int main(int argc, char *argv[])
{
    auto app = new QCoreApplication(argc, argv);
    app->setApplicationName("fimbeb-example");

    auto note = Note::newNote();
    note->set_title("yeet");
    note->save();

    auto child = Note::newNote();
    child->set_title("ohno");
    child->save();

    note->addChildNote(child);

    auto data = Note::where({ like(title, QStringLiteral("%ye%")) });

    auto model = new NoteModel;

    app->exec();

    delete model;

    return 0;
}

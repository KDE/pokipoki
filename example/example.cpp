#include <QCoreApplication>
#include "example.gen.h"

int main(int argc, char *argv[])
{
    auto app = new QCoreApplication(argc, argv);
    app->setApplicationName("fimbeb-example");

    auto note = Note::newNote();
    note->save();

    note->set_title("yeet one");
    note->save();

    note->set_title("yeet two");
    note->save();

    note->set_title("yeet three");
    note->save();

    for (; pUR->canUndo(); pUR->undo()) {
        qDebug() << note->title();
    }

    return 0;
}

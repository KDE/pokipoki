#include <QCoreApplication>
#include "example.gen.h"

int main(int argc, char *argv[])
{
    auto app = new QCoreApplication(argc, argv);
    app->setApplicationName("fimbeb-example");

    auto note = Note::newNote();
    Q_ASSERT(note);

    note->set_title("yeet one");
    note->save();

    note->set_title("yeet two");
    note->save();

    note->set_title("yeet three");
    note->save();

    for (; pUR->canUndo(); pUR->undo()) {
        qDebug() << note->title();
    }
    Q_ASSERT(note->title() == "yeet one");

    return 0;
}

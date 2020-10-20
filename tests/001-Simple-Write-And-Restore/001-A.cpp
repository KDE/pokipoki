#include <QCoreApplication>
#include "001.h"

int main(int argc, char* argv[]) {
    auto app = new QCoreApplication(argc, argv);
    app->setApplicationName("pokipoki-test-001");

    auto item = Item::newItem();
    item->set_prop("hi!");
    item->save();

    auto model = new ItemModel;
    if (model->rowCount() != 1) {
        return 1;
    }
    delete model;

    item->stageDelete();

    return 0;
}

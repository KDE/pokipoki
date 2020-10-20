#include <QCoreApplication>
#include "001.h"

int main(int argc, char* argv[]) {
    auto app = new QCoreApplication(argc, argv);
    app->setApplicationName("pokipoki-test-001");

    auto model = new ItemModel;
    if (model->rowCount() != 0) {
        return 1;
    }
    delete model;

    return 0;
}

-- +goose Up
-- +goose StatementBegin
INSERT INTO currency(id, code, rate) values ('a3c8a350-5b69-4d75-a16e-8d5bfa2b7a29','USD', 1);
INSERT INTO currency(id, code, rate) values ('bbd9c3f1-8a5f-4f3e-87e6-9c8b4a9d69c0','EUR', 1.5);
INSERT INTO currency(id, code, rate) values ('cc7a9d85-f728-4c44-b55b-34e354f5937a','RUB', 0.8);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM currency WHERE id = 'a3c8a350-5b69-4d75-a16e-8d5bfa2b7a29';
DELETE FROM currency WHERE id = 'bbd9c3f1-8a5f-4f3e-87e6-9c8b4a9d69c0';
DELETE FROM currency WHERE id = 'cc7a9d85-f728-4c44-b55b-34e354f5937a';
-- +goose StatementEnd

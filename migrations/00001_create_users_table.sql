-- +goose Up
-- +goose StatementBegin
CREATE TABLE
   users (
      id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
      name VARCHAR(255) NOT NULL,
      email VARCHAR(255) NOT NULL UNIQUE,
      age INT NOT NULL,
      balance DECIMAL NOT NULL
   );

INSERT INTO
   users (id, name, email, age, balance)
VALUES
   (
      'd290f1ee-6c54-4b01-90e6-d701748f0851',
      'A test',
      'a@example.com',
      25,
      100.50
   ),
   (
      'd290f1ee-6c54-4b01-90e6-d701748f0852',
      'B test',
      'b@example.com',
      30,
      75.25
   );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE users;

-- +goose StatementEnd
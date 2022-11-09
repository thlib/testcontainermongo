db.createCollection('myusers');
db.myusers.insertMany([
  {
      "id":    "9697f854-1f5e-48ef-a99a-000000000001",
      "email": "alice@example.com",
  },
  {
      "id":    "9697f854-1f5e-48ef-a99a-000000000002",
      "email": "bob@example.com",
  }
]);

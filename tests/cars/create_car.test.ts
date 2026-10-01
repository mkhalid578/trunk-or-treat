import assert from "node:assert/strict";
import { before, describe, it } from "node:test";
import request from "supertest";

const BASE_URL: string = "http://localhost:3000";
let authToken: string;

before(async () => {
  const email = `ci-${Date.now()}@example.com`;
  const password = "ci-test-password";

  const registerResponse = await request(BASE_URL)
    .post("/auth/register")
    .send({ email, password });
  assert.equal(registerResponse.status, 201);

  const loginResponse = await request(BASE_URL)
    .post("/auth/login")
    .send({ email, password });
  assert.equal(loginResponse.status, 200);
  assert.equal(typeof loginResponse.body.token, "string");
  authToken = loginResponse.body.token;
});

describe("POST /cars", () => {
  it("should return 400 for bad request", async () => {
    const response = await request(BASE_URL).post("/cars").set("Authorization", `Bearer ${authToken}`).send({
      make: "",
      model: "",
      model_year: 0,
      trim: "",
      body_style: "",
      powertrain: "",
    });
    assert.equal(response.status, 400);
  });

  it("should return 400 for missing required fields", async () => {
    const response = await request(BASE_URL).post("/cars").set("Authorization", `Bearer ${authToken}`).send({});
    assert.equal(response.status, 400);
  });

  it("should return 400 for invalid data types", async () => {
    const response = await request(BASE_URL).post("/cars").set("Authorization", `Bearer ${authToken}`).send({
      make: 123,
      model: 456,
      model_year: "invalid",
      trim: 789,
      body_style: 101112,
      powertrain: 131415,
    });
    assert.equal(response.status, 400);
  });

  it("should return 400 for negative model year", async () => {
    const response = await request(BASE_URL).post("/cars").set("Authorization", `Bearer ${authToken}`).send({
      make: "toyota",
      model: "rav4",
      model_year: -1,
      trim: "xle premium",
      body_style: "suv",
      powertrain: "hybrid",
    });
    assert.equal(response.status, 400);
  });

  it("should return 400 for model year in the future", async () => {
    const response = await request(BASE_URL).post("/cars").set("Authorization", `Bearer ${authToken}`).send({
      make: "toyota",
      model: "rav4",
      model_year: 3000,
      trim: "xle premium",
      body_style: "suv",
      powertrain: "hybrid",
    });
    assert.equal(response.status, 400);
  });

  it("should return 400 for invalid powertrain", async () => {
    const response = await request(BASE_URL).post("/cars").set("Authorization", `Bearer ${authToken}`).send({
      make: "toyota",
      model: "rav4",
      model_year: 2020,
      trim: "xle premium",
      body_style: "suv",
      powertrain: "invalid",
    });
    assert.equal(response.status, 400);
  });
  it("should return 400 for invalid body style", async () => {
    const response = await request(BASE_URL).post("/cars").set("Authorization", `Bearer ${authToken}`).send({
      make: "toyota",
      model: "rav4",
      model_year: 2020,
      trim: "xle premium",
      body_style: "invalid",
      powertrain: "hybrid",
    });
    assert.equal(response.status, 400);
  });

  it("should return 201 for valid car creation and 200 for deletion", async () => {
    const response = await request(BASE_URL).post("/cars").set("Authorization", `Bearer ${authToken}`).send({
      make: "toyota",
      model: "rav4",
      model_year: 2020,
      trim: "xle premium",
      body_style: "suv",
      powertrain: "hybrid",
    });
    assert.equal(response.status, 201);
    const createdCarId = response.body.id;

    const deleteResponse = await request(BASE_URL).delete(`/cars/${createdCarId}`).set("Authorization", `Bearer ${authToken}`);
    assert.equal(deleteResponse.status, 200);
  });
  
});

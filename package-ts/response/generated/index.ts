import { makeApi, Zodios, type ZodiosOptions } from "@zodios/core";
import { z } from "zod";

const model_HealthCheck = z.object({ status: z.string() }).passthrough();
const common_Error = z
  .object({
    type: z.string().url().optional(),
    title: z.string(),
    status: z.number().int(),
    detail: z.string().optional(),
    instance: z.string().optional(),
    traceId: z.string().optional(),
    errors: z
      .array(
        z
          .object({ field: z.string(), message: z.string() })
          .partial()
          .passthrough()
      )
      .optional(),
  })
  .passthrough();
const model_System = z
  .object({
    id: z.string().uuid(),
    systemName: z.string(),
    localGovernmentId: z.string().nullish(),
    createdAt: z.string().datetime({ offset: true }),
    updatedAt: z.string().datetime({ offset: true }),
    mailAddress: z.string().email(),
    telephone: z.string().nullish(),
    remark: z.string().nullish(),
  })
  .passthrough();

export const schemas = {
  model_HealthCheck,
  common_Error,
  model_System,
};

const endpoints = makeApi([
  {
    method: "get",
    path: "/api/v1/systems",
    alias: "GetSystems",
    description: `Retrieve a list of all systems with optional search filters`,
    requestFormat: "json",
    parameters: [
      {
        name: "systemName",
        type: "Query",
        schema: z.string().optional(),
      },
      {
        name: "email",
        type: "Query",
        schema: z.string().email().optional(),
      },
      {
        name: "localGovernmentId",
        type: "Query",
        schema: z.string().optional(),
      },
    ],
    response: z.array(model_System),
    errors: [
      {
        status: 500,
        description: `Internal Server Error`,
        schema: common_Error,
      },
    ],
  },
  {
    method: "post",
    path: "/api/v1/systems",
    alias: "CreateSystem",
    description: `Create a new system`,
    requestFormat: "json",
    parameters: [
      {
        name: "body",
        type: "Body",
        schema: model_System,
      },
    ],
    response: model_System,
    errors: [
      {
        status: 400,
        description: `Bad Request`,
        schema: common_Error,
      },
      {
        status: 500,
        description: `Internal Server Error`,
        schema: common_Error,
      },
    ],
  },
  {
    method: "get",
    path: "/api/v1/systems/:id",
    alias: "GetSystemById",
    description: `Retrieve a specific system by its ID`,
    requestFormat: "json",
    parameters: [
      {
        name: "id",
        type: "Path",
        schema: z.string().uuid(),
      },
    ],
    response: model_System,
    errors: [
      {
        status: 404,
        description: `System not found`,
        schema: common_Error,
      },
      {
        status: 500,
        description: `Internal Server Error`,
        schema: common_Error,
      },
    ],
  },
  {
    method: "put",
    path: "/api/v1/systems/:id",
    alias: "UpdateSystem",
    description: `Update an existing system`,
    requestFormat: "json",
    parameters: [
      {
        name: "body",
        type: "Body",
        schema: model_System,
      },
      {
        name: "id",
        type: "Path",
        schema: z.string().uuid(),
      },
    ],
    response: model_System,
    errors: [
      {
        status: 400,
        description: `Bad Request`,
        schema: common_Error,
      },
      {
        status: 404,
        description: `System not found`,
        schema: common_Error,
      },
      {
        status: 500,
        description: `Internal Server Error`,
        schema: common_Error,
      },
    ],
  },
  {
    method: "delete",
    path: "/api/v1/systems/:id",
    alias: "DeleteSystem",
    description: `Delete an existing system`,
    requestFormat: "json",
    parameters: [
      {
        name: "id",
        type: "Path",
        schema: z.string().uuid(),
      },
    ],
    response: z.void(),
    errors: [
      {
        status: 404,
        description: `System not found`,
        schema: common_Error,
      },
      {
        status: 500,
        description: `Internal Server Error`,
        schema: common_Error,
      },
    ],
  },
  {
    method: "get",
    path: "/health",
    alias: "HealthCheck",
    description: `Get the status of the server.`,
    requestFormat: "json",
    response: z.object({ status: z.string() }).passthrough(),
    errors: [
      {
        status: 500,
        description: `Server error`,
        schema: common_Error,
      },
    ],
  },
]);

export const api = new Zodios(endpoints);

export function createApiClient(baseUrl: string, options?: ZodiosOptions) {
  return new Zodios(baseUrl, endpoints, options);
}

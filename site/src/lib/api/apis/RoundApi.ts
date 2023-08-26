/* tslint:disable */
/* eslint-disable */
/**
 * Live Leaderboards
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * The version of the OpenAPI document: 0.1.0
 *
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */

import * as runtime from "../runtime";
import type { CreatedRound, Round, RoundRequest } from "../models/index";
import {
  CreatedRoundFromJSON,
  CreatedRoundToJSON,
  RoundFromJSON,
  RoundToJSON,
  RoundRequestFromJSON,
  RoundRequestToJSON,
} from "../models/index";

export interface RoundPostRequest {
  roundRequest?: RoundRequest;
}

export interface RoundRoundIDGetRequest {
  roundID: any;
}

export interface RoundRoundIDScorePutRequest {
  roundID: any;
  body?: any | null;
}

/**
 *
 */
export class RoundApi extends runtime.BaseAPI {
  /**
   * Creates a new round
   */
  async roundPostRaw(
    requestParameters: RoundPostRequest,
    initOverrides?: RequestInit | runtime.InitOverrideFunction,
  ): Promise<runtime.ApiResponse<CreatedRound>> {
    const queryParameters: any = {};

    const headerParameters: runtime.HTTPHeaders = {};

    headerParameters["Content-Type"] = "application/json";

    const response = await this.request(
      {
        path: `/round`,
        method: "POST",
        headers: headerParameters,
        query: queryParameters,
        body: RoundRequestToJSON(requestParameters.roundRequest),
      },
      initOverrides,
    );

    return new runtime.JSONApiResponse(response, (jsonValue) =>
      CreatedRoundFromJSON(jsonValue),
    );
  }

  /**
   * Creates a new round
   */
  async roundPost(
    requestParameters: RoundPostRequest = {},
    initOverrides?: RequestInit | runtime.InitOverrideFunction,
  ): Promise<CreatedRound> {
    const response = await this.roundPostRaw(requestParameters, initOverrides);
    return await response.value();
  }

  /**
   * Gets current round information
   */
  async roundRoundIDGetRaw(
    requestParameters: RoundRoundIDGetRequest,
    initOverrides?: RequestInit | runtime.InitOverrideFunction,
  ): Promise<runtime.ApiResponse<Round>> {
    if (
      requestParameters.roundID === null ||
      requestParameters.roundID === undefined
    ) {
      throw new runtime.RequiredError(
        "roundID",
        "Required parameter requestParameters.roundID was null or undefined when calling roundRoundIDGet.",
      );
    }

    const queryParameters: any = {};

    const headerParameters: runtime.HTTPHeaders = {};

    const response = await this.request(
      {
        path: `/round/{roundID}`.replace(
          `{${"roundID"}}`,
          encodeURIComponent(String(requestParameters.roundID)),
        ),
        method: "GET",
        headers: headerParameters,
        query: queryParameters,
      },
      initOverrides,
    );

    return new runtime.JSONApiResponse(response, (jsonValue) =>
      RoundFromJSON(jsonValue),
    );
  }

  /**
   * Gets current round information
   */
  async roundRoundIDGet(
    requestParameters: RoundRoundIDGetRequest,
    initOverrides?: RequestInit | runtime.InitOverrideFunction,
  ): Promise<Round> {
    const response = await this.roundRoundIDGetRaw(
      requestParameters,
      initOverrides,
    );
    return await response.value();
  }

  /**
   * Send a score event
   */
  async roundRoundIDScorePutRaw(
    requestParameters: RoundRoundIDScorePutRequest,
    initOverrides?: RequestInit | runtime.InitOverrideFunction,
  ): Promise<runtime.ApiResponse<void>> {
    if (
      requestParameters.roundID === null ||
      requestParameters.roundID === undefined
    ) {
      throw new runtime.RequiredError(
        "roundID",
        "Required parameter requestParameters.roundID was null or undefined when calling roundRoundIDScorePut.",
      );
    }

    const queryParameters: any = {};

    const headerParameters: runtime.HTTPHeaders = {};

    headerParameters["Content-Type"] = "application/json";

    const response = await this.request(
      {
        path: `/round/{roundID}/score`.replace(
          `{${"roundID"}}`,
          encodeURIComponent(String(requestParameters.roundID)),
        ),
        method: "PUT",
        headers: headerParameters,
        query: queryParameters,
        body: requestParameters.body as any,
      },
      initOverrides,
    );

    return new runtime.VoidApiResponse(response);
  }

  /**
   * Send a score event
   */
  async roundRoundIDScorePut(
    requestParameters: RoundRoundIDScorePutRequest,
    initOverrides?: RequestInit | runtime.InitOverrideFunction,
  ): Promise<void> {
    await this.roundRoundIDScorePutRaw(requestParameters, initOverrides);
  }
}
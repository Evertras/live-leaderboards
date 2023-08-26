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

import { exists, mapValues } from "../runtime";
/**
 *
 * @export
 * @interface PlayerData
 */
export interface PlayerData {
  /**
   *
   * @type {any}
   * @memberof PlayerData
   */
  name: any | null;
}

/**
 * Check if a given object implements the PlayerData interface.
 */
export function instanceOfPlayerData(value: object): boolean {
  let isInstance = true;
  isInstance = isInstance && "name" in value;

  return isInstance;
}

export function PlayerDataFromJSON(json: any): PlayerData {
  return PlayerDataFromJSONTyped(json, false);
}

export function PlayerDataFromJSONTyped(
  json: any,
  ignoreDiscriminator: boolean,
): PlayerData {
  if (json === undefined || json === null) {
    return json;
  }
  return {
    name: json["name"],
  };
}

export function PlayerDataToJSON(value?: PlayerData | null): any {
  if (value === undefined) {
    return undefined;
  }
  if (value === null) {
    return null;
  }
  return {
    name: value.name,
  };
}

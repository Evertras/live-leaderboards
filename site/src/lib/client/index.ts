import createClient from "openapi-fetch";
import { components, paths } from "./api";

const { GET, POST, PUT } = createClient<paths>({
  baseUrl:
    process.env.REACT_APP_EVERTRAS_API_BASE_PATH ||
    "https://leaderboard-api.evertras.com",
});

export type Course = components["schemas"]["Course"];
export type Round = components["schemas"]["Round"];
export type RoundID = components["schemas"]["RoundID"];
export type RoundRequest = components["schemas"]["RoundRequest"];
export type CreatedRound = components["schemas"]["CreatedRound"];

export async function getRoundByID(id: string): Promise<Round> {
  const { data, error } = await GET("/round/{roundID}", {
    params: {
      path: {
        roundID: id,
      },
    },
  });

  if (data) {
    return data;
  }

  throw error;
}

export async function createRound(
  request: RoundRequest,
): Promise<CreatedRound> {
  const { data, error } = await POST("/round", {
    body: request,
  });

  if (data) {
    return data;
  }

  throw error;
}

export async function getLatestRoundID(): Promise<RoundID> {
  const { data, error } = await GET("/latest/roundID", {});

  if (data) {
    return data;
  }

  throw error;
}

export async function submitScore(
  roundID: string,
  playerIndex: number,
  hole: number,
  score: number,
): Promise<void> {
  const { error } = await PUT("/round/{roundID}/score", {
    params: {
      path: {
        roundID: roundID,
      },
    },
    body: {
      playerIndex,
      hole,
      score,
    },
  });

  if (error) {
    throw error;
  }
}

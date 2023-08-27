import {
  Configuration,
  Round,
  RoundApi,
  RoundRequest,
  RoundRoundIDScorePutRequest,
} from "../api";

const configuration = new Configuration({
  basePath:
    process.env.REACT_APP_EVERTRAS_API_BASE_PATH ||
    "https://leaderboard-api.evertras.com",
});

const roundAPI = new RoundApi(configuration);

export function getRoundByID(id: string): Promise<Round> {
  return roundAPI.roundRoundIDGet({
    roundID: id,
  });
}

export async function createRound(request: RoundRequest): Promise<string> {
  return roundAPI.roundPost({ roundRequest: request }).then((res) => res.id);
}

export async function getLatestRoundID(): Promise<string> {
  const id: string = await roundAPI.latestRoundIDGet();

  return id;
}

export async function submitScore(
  roundID: string,
  playerIndex: number,
  hole: number,
  score: number,
): Promise<void> {
  const reqParams: RoundRoundIDScorePutRequest = {
    roundID: roundID,
    body: {
      playerIndex,
      hole,
      score,
    },
  };

  return roundAPI.roundRoundIDScorePut(reqParams);
}

export { roundAPI };

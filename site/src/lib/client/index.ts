import { Configuration, Round, RoundApi, RoundRequest } from "../api";

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

export { roundAPI };

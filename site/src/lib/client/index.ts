import { Configuration, Round, RoundApi } from "../api";

const configuration = new Configuration({
  basePath:
    process.env.REACT_APP_EVERTRAS_API_BASE_PATH ?? "http://localhost:8037",
});

const roundAPI = new RoundApi(configuration);

export function getRoundByID(id: string): Promise<Round> {
  return roundAPI.roundRoundIDGet({
    roundID: id,
  });
}

export { roundAPI };

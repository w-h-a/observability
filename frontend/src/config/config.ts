export class Config {
	private static Instance: Config;
	private configData: Record<string, string>;

	private constructor() {
		// TODO: figure this out
		this.configData = {
			baseUrl: "http://localhost:4000/api/v1",
		};
	}

	public static GetInstance(): Config {
		if (!Config.Instance) {
			Config.Instance = new Config();
		}

		return Config.Instance;
	}

	public get(key: string): string {
		return this.configData[key];
	}
}

export class Config {
	private static Instance: Config;
	private configData: Record<string, string>;

	private constructor() {
		const envVars = process.env;

		this.configData = {
			BASE_QUERY_URL: "http://localhost:4000/api/v1",
			ENVIRONMENT: "dev",
		};

		if (envVars["BASE_QUERY_URL"]) {
			this.configData.BASE_QUERY_URL = envVars["BASE_QUERY_URL"];
		}

		if (envVars["NODE_ENV"]) {
			this.configData.ENVIRONMENT = envVars["NODE_ENV"];
		}
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

export enum EnvVar {
	ENVIRONMENT = "ENVIRONMENT",
	BASE_QUERY_URL = "BASE_QUERY_URL",
}

export class Config {
	private static Instance: Config;
	private configData: Record<EnvVar, string>;

	private constructor() {
		const envVars = process.env;

		this.configData = {
			[EnvVar.ENVIRONMENT]: "dev",
			[EnvVar.BASE_QUERY_URL]: "http://localhost:4000/api/v1",
		};

		if (envVars["NODE_ENV"]) {
			this.configData.ENVIRONMENT = envVars["NODE_ENV"];
		}

		if (envVars["BASE_QUERY_URL"]) {
			this.configData.BASE_QUERY_URL = envVars["BASE_QUERY_URL"];
		}
	}

	public static GetInstance(): Config {
		if (!Config.Instance) {
			Config.Instance = new Config();
		}

		return Config.Instance;
	}

	public get(key: EnvVar): string {
		return this.configData[key];
	}
}

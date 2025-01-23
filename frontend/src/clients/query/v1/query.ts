import { Config, EnvVar } from "../../../config/config";
import { IClient } from "../client";
import { FilteredQuery } from "../filteredQuery";

export class Query {
	static async GetServices<T = any>(
		client: IClient,
		start: number,
		end: number,
	): Promise<{ data: T }> {
		try {
			const path = `/services`;

			const query = `?start=${start}&end=${end}`;

			const rsp = await client.get<T>(
				`${Config.GetInstance().get(EnvVar.BASE_QUERY_URL)}${path}${query}`,
			);

			return rsp;
		} catch (err: unknown) {
			console.log(`query client failed to retrieve services: ${err}`);
			throw err;
		}
	}

	static async GetServiceNames<T = any>(client: IClient): Promise<{ data: T }> {
		try {
			const path = `/services/list`;

			const rsp = await client.get<T>(
				`${Config.GetInstance().get(EnvVar.BASE_QUERY_URL)}${path}`,
			);

			return rsp;
		} catch (err: unknown) {
			console.log(`query client failed to retrieve service names: ${err}`);
			throw err;
		}
	}

	static async GetOperationNames<T = any>(
		client: IClient,
		serviceName: string,
	): Promise<{ data: T }> {
		try {
			const path = `/service/operations`;

			const query = `?service=${serviceName}`;

			const rsp = await client.get<T>(
				`${Config.GetInstance().get(EnvVar.BASE_QUERY_URL)}${path}${query}`,
			);

			return rsp;
		} catch (err: unknown) {
			console.log(`query client failed to retrieve operation names: ${err}`);
			throw err;
		}
	}

	static async GetTags<T = any>(
		client: IClient,
		serviceName: string,
	): Promise<{ data: T }> {
		try {
			const path = `/service/tags`;

			const query = `?service=${serviceName}`;

			const rsp = await client.get<T>(
				`${Config.GetInstance().get(EnvVar.BASE_QUERY_URL)}${path}${query}`,
			);

			return rsp;
		} catch (err: unknown) {
			console.log(`query client failed to retrieve tags: ${err}`);
			throw err;
		}
	}

	static async GetEndpoints<T = any>(
		client: IClient,
		start: number,
		end: number,
		serviceName: string,
	): Promise<{ data: T }> {
		try {
			const path = `/service/endpoints`;

			const query = `?start=${start}&end=${end}&service=${serviceName}`;

			const rsp = await client.get<T>(
				`${Config.GetInstance().get(EnvVar.BASE_QUERY_URL)}${path}${query}`,
			);

			return rsp;
		} catch (err: unknown) {
			console.log(`query client failed to retrieve endpoints: ${err}`);
			throw err;
		}
	}

	static async GetServiceMetrics<T = any>(
		client: IClient,
		start: number,
		end: number,
		serviceName: string,
	): Promise<{ data: T }> {
		try {
			const path = `/service/overview`;

			const query = `?start=${start}&end=${end}&step=60&service=${serviceName}`;

			const rsp = await client.get<T>(
				`${Config.GetInstance().get(EnvVar.BASE_QUERY_URL)}${path}${query}`,
			);

			return rsp;
		} catch (err: unknown) {
			console.log(`query client failed to retrieve service metrics: ${err}`);
			throw err;
		}
	}

	static async GetSpans<T = any>(
		client: IClient,
		start: number,
		end: number,
		filters?: FilteredQuery,
	): Promise<{ data: T }> {
		try {
			const path = `/spans`;

			let query = `?start=${start}&end=${end}`;

			if (filters) {
				if (filters.service) {
					query += `&service=${filters.service}`;
				}

				if (filters.operation) {
					query += `&name=${filters.operation}`;
				}

				if (filters.duration) {
					if (filters.duration.min) {
						query += `&minDuration=${filters.duration.min}`;
					}

					if (filters.duration.max) {
						query += `&maxDuration=${filters.duration.max}`;
					}
				}

				if (filters.tags.length !== 0) {
					query += `&tags=${encodeURIComponent(JSON.stringify(filters.tags))}`;
				}
			}

			const rsp = await client.get<T>(
				`${Config.GetInstance().get(EnvVar.BASE_QUERY_URL)}${path}${query}`,
			);

			return rsp;
		} catch (err: unknown) {
			console.log(`query client failed to retrieve spans: ${err}`);
			throw err;
		}
	}

	static async GetSpansByTraceId<T = any>(
		client: IClient,
		traceId: string,
	): Promise<{ data: T }> {
		try {
			const path = `/spans/trace`;

			const query = `?traceId=${traceId}`;

			const rsp = await client.get<T>(
				`${Config.GetInstance().get(EnvVar.BASE_QUERY_URL)}${path}${query}`,
			);

			return rsp;
		} catch (err: unknown) {
			console.log(`query client failed to retrieve spans by trace id: ${err}`);
			throw err;
		}
	}
}

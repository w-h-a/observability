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

	static async GetServiceDependencies<T = any>(
		client: IClient,
		start: number,
		end: number,
	): Promise<{ data: T }> {
		try {
			const path = `/services/dependencies`;

			const query = `?start=${start}&end=${end}`;

			const rsp = await client.get<T>(
				`${Config.GetInstance().get(EnvVar.BASE_QUERY_URL)}${path}${query}`,
			);

			return rsp;
		} catch (err: unknown) {
			console.log(`query client failed to retrieve service dependencies: ${err}`);
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

	static async GetTraces<T = any>(
		client: IClient,
		start: number,
		end: number,
		service: string,
		traceId: string,
	): Promise<{ data: T }> {
		try {
			const path = `/traces`;

			let query = `?start=${start}&end=${end}`;

			if (service.length !== 0) {
				query += `&service=${service}`;
			}

			if (traceId.length !== 0) {
				query += `&traceId=${traceId}`;
			}

			const rsp = await client.get<T>(
				`${Config.GetInstance().get(EnvVar.BASE_QUERY_URL)}${path}${query}`,
			);

			return rsp;
		} catch (err: unknown) {
			console.log(`query client failed to retrieve traces: ${err}`);
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
				query = Query.applyFilters(filters, query);
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

	static async GetAggregation<T = any>(
		client: IClient,
		start: number,
		end: number,
		dimension: string,
		aggregation: string,
		filters?: FilteredQuery,
	): Promise<{ data: T }> {
		try {
			const path = `/spans/aggregated`;

			let query = `?start=${start}&end=${end}&dimension=${dimension}&aggregation=${aggregation}&step=60`;

			if (filters) {
				query = Query.applyFilters(filters, query);
			}

			const rsp = await client.get<T>(
				`${Config.GetInstance().get(EnvVar.BASE_QUERY_URL)}${path}${query}`,
			);

			return rsp;
		} catch (err: unknown) {
			console.log(`query client failed to retrieve aggregations: ${err}`);
			throw err;
		}
	}

	static async GetMetrics<T = any>(
		client: IClient,
		start: number,
		end: number,
		dimension: string,
		aggregation: string,
		filters?: FilteredQuery,
	): Promise<{ data: T }> {
		try {
			const path = `/metrics`;

			let query = `?start=${start}&end=${end}&dimension=${dimension}&aggregation=${aggregation}&step=60`;

			if (filters) {
				query = Query.applyFilters(filters, query);
			}

			const rsp = await client.get<T>(
				`${Config.GetInstance().get(EnvVar.BASE_QUERY_URL)}${path}${query}`,
			);

			return rsp;
		} catch (err: unknown) {
			console.log(`query client failed to retrieve metrics: ${err}`);
			throw err;
		}
	}

	private static applyFilters(filters: FilteredQuery, query: string) {
		if (filters.service) {
			query += `&service=${filters.service}`;
		}

		if (filters.operation) {
			query += `&name=${filters.operation}`;
		}

		if (filters.kind) {
			query += `&kind=${filters.kind}`;
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

		return query;
	}
}

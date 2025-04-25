// to use with @duplojs/http-client
/* eslint-disable */
/* prettier-ignore */
/* istanbul ignore file */
/* v8 ignore start */
// noinspection JSUnusedGlobalSymbols
// @ts-nocheck

type kvltRoutes = ({
    method: "PUT";
    path: "/value";
    body: {
        key: string;
        value: any;
        duration: number;
    }
    response: {
        code: 400;
        information: "error.keyRequired" | "error.bodyParams";
        body: {
            error: string;
        }
    } | {
        code: 200;
        information: "success.keySet";
    }
}) | ({
    method: "GET";
    path: "/value/{key}";
    params: {
        key: string;
    };
    response: {
        code: 400;
        information: "error.routeParams";
        body: {
            error: string;
        }
    } | {
        code: 404;
        information: "error.keyNotFound";
        body: {
            error: string;
        }
    } | {
        code: 200;
        information: "success.keyFound";
        body: {
            data: {
                key: string;
                value: any;
            }
        }
    }
});

export { kvltRoutes };
/* v8 ignore stop */
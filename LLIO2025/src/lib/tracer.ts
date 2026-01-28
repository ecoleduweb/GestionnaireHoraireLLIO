import { WebTracerProvider } from '@opentelemetry/sdk-trace-web';
import { getWebAutoInstrumentations } from '@opentelemetry/auto-instrumentations-web';
import { SimpleSpanProcessor } from '@opentelemetry/sdk-trace-base';
import { registerInstrumentations } from '@opentelemetry/instrumentation';
import { ZoneContextManager } from '@opentelemetry/context-zone';
import { OTLPTraceExporter } from '@opentelemetry/exporter-trace-otlp-http';
import { W3CTraceContextPropagator } from "@opentelemetry/core";
import { ATTR_SERVICE_NAME } from '@opentelemetry/semantic-conventions';
import { resourceFromAttributes } from '@opentelemetry/resources';
import { env } from "$env/dynamic/public";

const TRACE_URL = env.PUBLIC_TRACE_URL;
const APPLICATION_NAME = env.PUBLIC_APPLICATION_NAME;

const exporter = new OTLPTraceExporter({
    url: TRACE_URL,
    headers: {},
});

const provider = new WebTracerProvider({
    resource: resourceFromAttributes({
        [ATTR_SERVICE_NAME]: APPLICATION_NAME,
    }),
    spanProcessors: [new SimpleSpanProcessor(exporter)],
});

provider.register({
    contextManager: new ZoneContextManager(),
    propagator: new W3CTraceContextPropagator(),
});


export class ClientTelemetry {
    private static instance: ClientTelemetry;
    private initialized = false;
    private tracer;

    private constructor() {
        this.tracer = provider.getTracer('web-service-tracer');
    }

    public static getInstance(): ClientTelemetry {
        if (!ClientTelemetry.instance) {
            ClientTelemetry.instance = new ClientTelemetry();
        }
        return ClientTelemetry.instance;
    }

    public start(): void {
        if (!this.initialized) {
            registerInstrumentations({
                instrumentations: [
                    getWebAutoInstrumentations({
                        // You can configure specific instrumentations here
                        '@opentelemetry/instrumentation-fetch': {
                            propagateTraceHeaderCorsUrls: [
                                /.+/g, // Propagate to all URLs - customize this according to your needs
                            ],
                        },
                    }),
                ],
            });
            console.log("Client Telemetry Initialized");
            this.initialized = true;
        }
    }

    // Utility method to create spans
    public async trackSpan<T>(name: string, fn: () => Promise<T>): Promise<T> {
        const span = this.tracer.startSpan(name);
        try {
            const result = await fn();
            span.setStatus({ code: 0 });
            return result;
        } catch (error: any) {
            span.setStatus({ code: 1, message: error.message });
            span.recordException(error);
            throw error;
        } finally {
            span.end();
        }
    }
}

// Usage example:
/*
// In your main.ts or app initialization
const telemetry = ClientTelemetry.getInstance();
telemetry.start();

// In your components
async function handleClick() {
    const telemetry = ClientTelemetry.getInstance();
    await telemetry.trackSpan('button-click', async () => {
        const response = await fetch('/api/data');
        return response.json();
    });
}
*/
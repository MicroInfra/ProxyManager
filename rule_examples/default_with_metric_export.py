from mitmproxy import proxy, options
from mitmproxy.tools.dump import DumpMaster
from mitmproxy.addons import core
from mitmproxy.http import HTTPFlow
from threading import Thread
import time
from urllib.parse import urlparse
import re
from os import getenv


SERVICE_NAME = '{{service_name}}'
AUTH_TOKEN = '{{exporter_auth_token}}'
EXPORTER_URL = '{{exporter_url}}'


def filter_rules(flow: HTTPFlow):
    metrics = []

    # # Examples of filtering
    # if '\' or 1=1;--' in flow.request.url or \
    #         'nginx' in flow.response.get_text() or\
    #         flow.request.pretty_url.endswith("/something"):
    #     tmp = flow.response.text.replace('nginx', 'mark_was_here')
    #     metrics.append(format_metric('filtered_request_rule1', 1, 'sum'))
    #     flow.response.text = tmp # or even  = 'some text'
    #     flow.response.status_code = 418

    return metrics


class Addon:

    def __init__(self):
        pass
    def response(self, flow: HTTPFlow):
        start = time.time()
        specific_metrics = filter_rules(flow)
        thread = Thread(target=send_metrics, args=(flow, start, specific_metrics))
        thread.start()


def send_metrics(flow: HTTPFlow, start_time, specific_metrics: list or None = None):
    import requests
    metrics = []
    if specific_metrics:
        metrics = specific_metrics

    metrics.append(format_metric('request_count', 1, 'counter'))
    request_timing = (flow.request.timestamp_end - flow.request.timestamp_start) * 1000
    response_timing = (flow.response.timestamp_end - flow.response.timestamp_start) * 1000
    metrics.append(format_metric('request_latency_ms', request_timing))
    metrics.append(format_metric('response_latency_ms', response_timing))

    metrics.append(format_metric('response_code', 1, 'counter', {"code": flow.response.status_code}))

    path = urlparse(flow.request.pretty_url).path
    metrics.append(format_metric('url_rps', 1, 'counter', {"url": path}))

    metrics.append(format_metric('response_size', len(flow.response.get_text()), 'sum'))
    metrics.append(format_metric('request_size', len(flow.request.get_text()), 'sum'))

    # Maybe add request ip. But need to do some modifications in exporter

    proxy_ms = (time.time() - start_time) * 1000
    metrics.append(format_metric('proxy_latency_ms', proxy_ms))

    requests.post(EXPORTER_URL, json={'metrics': metrics}, headers={"Authorization": AUTH_TOKEN}, timeout=1)


def format_metric(name, value, type='gauge', labels=None):
    if labels:
        return {'name': SERVICE_NAME + name, 'value': value, 'type': type, 'labels': labels}
    return {'name': SERVICE_NAME + name, 'value': value, 'type': type}


addons = [
    Addon()
]

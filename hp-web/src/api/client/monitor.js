import http from '../../data/http';

export function monitorList(query) {
  return http({
    url: '/client/monitor/list',
    method: 'get',
    params: query,
  });
}

export function getMonitorData(query) {
  return http({
    url: '/client/monitor/getMonitorData',
    method: 'get',
    params: query,
  });
}

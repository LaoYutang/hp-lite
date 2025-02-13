<template>
  <div>
    <div v-if="monitorData && Object.keys(monitorData).length === 0">暂无穿透数据统计</div>
    <div v-else id="chat" />
  </div>
</template>

<script setup>
  import * as echarts from 'echarts/core';
  import {
    TitleComponent,
    ToolboxComponent,
    TooltipComponent,
    GridComponent,
    LegendComponent,
    DataZoomComponent,
    MarkAreaComponent,
  } from 'echarts/components';
  import { LineChart } from 'echarts/charts';
  import { UniversalTransition } from 'echarts/features';
  import { CanvasRenderer } from 'echarts/renderers';
  import { onMounted, ref } from 'vue';
  import { monitorList } from '../../api/client/monitor';
  import { getConfigList } from '../../api/client/config';

  echarts.use([
    TitleComponent,
    ToolboxComponent,
    TooltipComponent,
    GridComponent,
    LegendComponent,
    DataZoomComponent,
    MarkAreaComponent,
    LineChart,
    CanvasRenderer,
    UniversalTransition,
  ]);

  const monitorData = ref();
  const currentConfigList = ref();

  const loadData = async () => {
    let data = await monitorList();
    monitorData.value = data.data;
  };

  const formatTimestamp = (timestamp) => {
    const date = new Date(timestamp);
    const year = date.getFullYear();
    const month = padZero(date.getMonth() + 1);
    const day = padZero(date.getDate());
    const hours = padZero(date.getHours());
    const minutes = padZero(date.getMinutes());
    const seconds = padZero(date.getSeconds());
    return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
  };

  const padZero = (num) => {
    return num.toString().padStart(2, '0');
  };

  const getRemarkById = (id) => {
    const remark = currentConfigList.value.find((item) => item.id.toString() === id)?.remarks ?? '';
    return remark;
  };

  const showFlow = (key, dataList) => {
    var chartDom = document.getElementById('chat');
    let flow = document.createElement('div');
    flow.id = 'flow' + key;
    flow.style.width = '100%';
    flow.style.height = '40vh';
    chartDom.appendChild(flow);

    let option = {
      title: {
        text: `配置 ${getRemarkById(key)} (id:${key}) 下载/上传`,
        left: 'center',
      },
      tooltip: {
        trigger: 'axis',
        axisPointer: {
          type: 'cross',
          animation: false,
          label: {
            backgroundColor: '#505765',
          },
        },
      },
      legend: {
        data: ['下载', '上传'],
        top: 30,
      },
      xAxis: [
        {
          type: 'category',
          boundaryGap: false,
          axisLine: { onZero: false },
          // prettier-ignore
          data: dataList.map(item => formatTimestamp(item.time)),
        },
      ],
      yAxis: [
        {
          name: '下载/MB',
          type: 'value',
        },
        {
          name: '上传/MB',
          type: 'value',
        },
      ],
      series: [
        {
          name: '下载',
          type: 'line',
          smooth: true,
          emphasis: {
            focus: 'series',
          },
          // prettier-ignore
          data: dataList.map(item => (item.download / 1024/1024).toFixed(2)),
        },
        {
          name: '上传',
          type: 'line',
          smooth: true,
          yAxisIndex: 1,
          emphasis: {
            focus: 'series',
          },
          // prettier-ignore
          data: dataList.map(item => (item.upload / 1024/1024).toFixed(2)),
        },
      ],
    };

    let myChart = echarts.init(document.getElementById(flow.id));
    option && myChart.setOption(option);
  };

  const showAccess = (key, dataList) => {
    var chartDom = document.getElementById('chat');
    let flow = document.createElement('div');
    flow.id = 'access' + key;
    flow.style.width = '100%';
    flow.style.height = '40vh';
    chartDom.appendChild(flow);

    let option = {
      title: {
        text: `配置 ${getRemarkById(key)} (id:${key}) pv/uv 统计`,
        left: 'center',
      },
      tooltip: {
        trigger: 'axis',
      },
      legend: {
        data: ['pv', 'uv'],
        top: 30,
      },
      xAxis: {
        type: 'category',
        boundaryGap: false,
        data: dataList.map((item) => formatTimestamp(item.time)),
      },
      yAxis: [
        {
          name: '访问量/人数',
          type: 'value',
        },
      ],

      series: [
        {
          name: 'pv',
          type: 'line',
          smooth: true,
          data: dataList.map((item) => item.pv),
        },
        {
          name: 'uv',
          type: 'line',
          smooth: true,
          data: dataList.map((item) => item.uv),
        },
      ],
    };

    let myChart = echarts.init(document.getElementById(flow.id));
    option && myChart.setOption(option);
  };

  const loadConfigData = async () => {
    currentConfigList.value = [];
    await getConfigList({
      total: 0,
      current: 1,
      pageSize: 1000,
    }).then((res) => {
      currentConfigList.value = res.data.records;
    });
  };

  onMounted(async () => {
    await loadData();
    await loadConfigData();
    for (let monitorDataKey in monitorData.value) {
      showFlow(monitorDataKey, monitorData.value[monitorDataKey]);
      showAccess(monitorDataKey, monitorData.value[monitorDataKey]);
    }
  });
</script>

<style scoped />

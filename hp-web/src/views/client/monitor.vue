<template>
  <div>
    <div v-if="currentConfigList && currentConfigList.length === 0">暂无配置</div>
    <div v-else>
      <a-select style="width: 200px" v-model:value="currentConfigId" @change="selectChange">
        <a-select-option v-for="item in currentConfigList" :key="item.id" :value="item.id">
          {{ `${item.remarks} (id:${item.id})` }}
        </a-select-option>
      </a-select>
      <div id="chart">
        <div id="flow" style="width: 100%; height: 40vh" />
        <div id="access" style="width: 100%; height: 40vh" />
      </div>
    </div>
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
  import { getMonitorData } from '../../api/client/monitor';
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

  const currentConfigList = ref();
  const currentConfigId = ref();

  const loadData = async (configId) => {
    let data = await getMonitorData({ configId });
    return data.data;
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
    var chartDom = document.getElementById('flow');

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
          data: dataList.map((item) => (item.download / 1024 / 1024).toFixed(2)),
        },
        {
          name: '上传',
          type: 'line',
          smooth: true,
          yAxisIndex: 1,
          emphasis: {
            focus: 'series',
          },
          data: dataList.map((item) => (item.upload / 1024 / 1024).toFixed(2)),
        },
      ],
    };

    let flowChart = echarts.init(chartDom);
    option && flowChart.setOption(option);
  };

  const showAccess = (key, dataList) => {
    var chartDom = document.getElementById('access');

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

    let accessChart = echarts.init(chartDom);
    option && accessChart.setOption(option);
  };

  const selectChange = async (value) => {
    const res = await loadData(value);
    showFlow(value, res);
    showAccess(value, res);
  };

  const loadConfigData = async () => {
    currentConfigList.value = [];
    await getConfigList({
      total: 0,
      current: 1,
      pageSize: 1000,
    }).then((res) => {
      currentConfigList.value = res.data.records;
      if (currentConfigList.value.length > 0) {
        currentConfigId.value = currentConfigList.value[0].id;
      }
    });
  };

  onMounted(async () => {
    await loadConfigData();
    if (currentConfigList.value.length > 0) {
      selectChange(currentConfigList.value[0].id);
    }
  });
</script>

<style scoped />

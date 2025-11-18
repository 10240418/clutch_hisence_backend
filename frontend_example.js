// 质量统计API使用示例

// API接口地址
const QUALITY_STATS_API = '/api/management/quality_stats'

// 调用质量统计API的函数
async function getQualityStats(startDate, endDate) {
    try {
        const response = await fetch(`${QUALITY_STATS_API}?startDate=${startDate}&endDate=${endDate}`, {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('token')}`,
                'Content-Type': 'application/json'
            }
        })

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`)
        }

        const result = await response.json()
        return result.data
    } catch (error) {
        console.error('获取质量统计数据失败:', error)
        throw error
    }
}

// 使用示例
async function loadQualityDashboard() {
    const startDate = '2024-01-01'
    const endDate = '2024-12-31'

    try {
        const stats = await getQualityStats(startDate, endDate)

        // 更新合格率饼图数据
        qualityRateOption.value.series[0].data = [
            {
                value: stats.qualityRate.qualityRate,
                name: '合格',
                itemStyle: { color: '#52c41a' }
            },
            {
                value: 100 - stats.qualityRate.qualityRate,
                name: '不合格',
                itemStyle: { color: '#ff4d4f' }
            }
        ]

        // 更新不良类型分布饼图数据
        defectTypeOption.value.series[0].data = stats.defectTypeDistribution.map((item, index) => ({
            value: item.rate,
            name: item.type,
            itemStyle: { color: getColorByIndex(index) }
        }))

        // 更新供应商不良趋势数据
        const supplierSeries = stats.supplierDefectTrend.map((supplier, index) => ({
            name: supplier.supplierName,
            type: 'line',
            data: supplier.dailyData.map(day => day.defectRate),
            itemStyle: { color: getColorByIndex(index) }
        }))

        manufacturerTrendOption.value.series = supplierSeries
        manufacturerTrendOption.value.legend.data = stats.supplierDefectTrend.map(s => s.supplierName)

        // 更新各类不良趋势数据
        const defectTrends = stats.defectTrendByType

        // 更新绝缘耐压不良趋势
        const insulationOption = generateLineChartOption(
            '绝缘耐压不良',
            '#1890ff',
            defectTrends.insulationData.map(d => d.count)
        )

        // 更新电阻不良趋势
        const resistanceOption = generateLineChartOption(
            '电阻不良',
            '#722ed1',
            defectTrends.resistanceData.map(d => d.count)
        )

        // 更新反电动势不良趋势
        const emfOption = generateLineChartOption(
            '反电动势不良',
            '#fa8c16',
            defectTrends.emfData.map(d => d.count)
        )

        // 更新外观不良趋势
        const appearanceOption = generateLineChartOption(
            '外观不良',
            '#52c41a',
            defectTrends.appearanceData.map(d => d.count)
        )

        // 更新噪音不良趋势
        const noiseOption = generateLineChartOption(
            '噪音不良',
            '#eb2f96',
            defectTrends.noiseData.map(d => d.count)
        )

        console.log('质量统计数据加载完成:', stats)

    } catch (error) {
        console.error('加载质量统计数据失败:', error)
    }
}

// 辅助函数：根据索引获取颜色
function getColorByIndex(index) {
    const colors = ['#1890ff', '#722ed1', '#fa8c16', '#52c41a', '#eb2f96', '#13c2c2', '#f5222d']
    return colors[index % colors.length]
}

// 响应式数据更新示例（使用Vue 3 Composition API）
import { ref, onMounted } from 'vue'

export function useQualityStats() {
    const loading = ref(false)
    const dateRange = ref(['2024-01-01', '2024-12-31'])

    // 更新图表数据的响应式函数
    const updateCharts = async () => {
        loading.value = true
        try {
            const [startDate, endDate] = dateRange.value
            const stats = await getQualityStats(startDate, endDate)

            // 在这里更新你的图表选项
            // 参考上面的代码更新各个图表的数据

        } catch (error) {
            console.error('更新图表数据失败:', error)
        } finally {
            loading.value = false
        }
    }

    // 组件挂载时加载数据
    onMounted(() => {
        updateCharts()
    })

    return {
        loading,
        dateRange,
        updateCharts
    }
}

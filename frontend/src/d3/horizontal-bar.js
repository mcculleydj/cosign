import * as d3 from 'd3'

function transformData({ member, memberMap }, cutoff) {
  return Object.entries(member.counts)
    .map(([id, count]) => ({ count, ...memberMap[id] }))
    .sort((m1, m2) => m2.count - m1.count)
    .slice(0, cutoff)
}

function drawBars(svg, data, params, callbacks) {
  const { width, height, margin } = params
  const { onMouseover, onMouseleave, onMouseup } = callbacks

  const x = d3
    .scaleLinear()
    .domain([0, d3.max(data, d => d.count)])
    .range([margin.left, width - margin.right])

  const y = d3
    .scaleBand()
    .domain(d3.range(data.length))
    .rangeRound([margin.top, height - margin.bottom])
    .padding(0.2)

  svg
    .append('g')
    .attr('fill', 'steelblue')
    .selectAll('rect')
    .data(data)
    .join('rect')
    .attr('x', x(0))
    .attr('y', (_, i) => y(i))
    .attr('width', d => x(d.count) - x(0))
    .attr('height', y.bandwidth())
    .style('cursor', 'pointer')
    .on('mouseover', function(_, d) {
      onMouseover(d)
      d3.select(this).classed('active', true)
    })
    .on('mouseleave', function() {
      onMouseleave()
      d3.select(this).classed('active', false)
    })
    .on('mouseup', function(_, d) {
      onMouseup(d)
    })

  const xAxis = g =>
    g.attr('transform', `translate(0, ${margin.top})`).call(d3.axisTop(x))

  const yAxis = g =>
    g
      .attr('transform', `translate(${margin.left}, 0)`)
      .call(d3.axisLeft(y).tickValues([]))

  svg.append('g').call(xAxis)
  svg.append('g').call(yAxis)
}

export function drawChart(selector, raw, params, callbacks) {
  const data = transformData(raw, params.cutoff)

  const svg = d3
    .select(selector)
    .append('svg')
    .attr('viewBox', [0, 0, params.width, params.height])
    .attr('preserveAspectRatio', 'xMinYMin meet')
    .classed('svg-content', true)

  drawBars(svg, data, params, callbacks)

  return () => {
    svg.remove()
  }
}

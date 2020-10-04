import * as d3 from 'd3'

export function transform(member, memberMap, threshold) {
  return Object.keys(member.counts)
    .map(id => +id)
    .map(id => ({ ...memberMap[id], count: member.counts[id] }))
    .filter(m => m.count >= threshold)
    .map((m, i) => ({ ...m, index: i }))
    .sort((m1, m2) => m1.count - m2.count)
}

export function drawChart(selector, member, memberMap, params, callbacks) {
  const data = transform(member, memberMap, params.threshold)

  const svg = d3
    .select(selector)
    .append('svg')
    .attr('viewBox', [0, 0, params.width, params.height])

  // addContext(svg, data, params)

  buildRadialBarChart(svg, data, params, callbacks)

  function remove() {
    svg.remove()
  }

  return remove
}

function arc(xScale, yScale, params) {
  return d3
    .arc()
    .innerRadius(params.innerRadius)
    .outerRadius(d => yScale(d.count))
    .startAngle(d => xScale(d.index))
    .endAngle(d => xScale(d.index) + xScale.bandwidth())
    .padAngle(0.01)
    .padRadius(params.innerRadius)
}

function xScale(data) {
  return d3
    .scaleBand()
    .domain(data.map(d => d.index))
    .range([0, 2 * Math.PI])
}

function yScale(data, params) {
  return d3
    .scaleLinear()
    .domain(d3.extent(data, d => d.count))
    .range([params.innerRadius, outerRadius(params)])
}

function outerRadius(params) {
  return Math.min(params.height, params.width) / 2
}

function colors(data) {
  return d3
    .scaleSequential()
    .domain(d3.extent(data, d => d.count))
    .interpolator(d3.interpolateCividis)
}

function buildRadialBarChart(svg, data, params, callbacks) {
  let name
  svg
    .append('g')
    .attr('transform', `translate(${params.width / 2}, ${params.height / 2})`)
    .selectAll('path')
    .data(data)
    .join(enter => {
      enter
        .append('path')
        .attr('fill', d => colors(data)(d.count))
        .attr('stroke', '1px')
        .attr('d', arc(xScale(data), yScale(data, params), params))
        .on('mouseover', function(_, d) {
          callbacks.setCosponsor(d)
          d3.select(this).classed('active', true)

          name = svg
            .append('text')
            .text(d.name)
            .attr('x', params.width / 2 + 10)
            .attr('y', params.height / 6)
        })
        .on('mouseleave', function() {
          callbacks.setCosponsor(null)
          d3.select(this).classed('active', false)

          name.remove()
        })
    })
}

// function addContext(svg, data, params) {
// svg
//   .append('text')
//   .attr('class', 'title')
//   .text('Cosponsors')
//   .attr('x', params.width / 2 + 10)
//   .attr('y', params.height / 6)
// svg
//   .append('text')
//   .attr('class', 'content')
//   .text('Blurb')
//   .attr('x', params.width / 2 + 10)
//   .attr('y', params.height / 6 + 20)
// svg
//   .append('text')
//   .attr('font-size', params.bodyFontSize)
//   .text('Score')
//   .attr('x', params.width / 2 - params.innerRadius / 2)
//   .attr('y', params.height / 2 + params.innerRadius / 6)
// }

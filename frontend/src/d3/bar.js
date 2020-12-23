import * as d3 from 'd3'

export default class Bar {
  constructor(selector, margins) {
    this.selector = selector
    this.margins = margins

    this.svg = d3.select(`#${selector}`).append('svg')
    this.plot = this.svg
      .append('g')
      .attr('transform', `translate(${margins.left}, ${margins.top})`)

    // scales
    this.x = d3.scaleLinear()
    this.y = d3
      .scaleBand()
      .paddingInner(0.3)
      .paddingOuter(0.3)

    // y-axis
    this.yAxisFn = d3.axisLeft().tickValues([])
    this.yAxis = this.plot.append('g')
  }

  // data shape:
  // [{ [key]: string, value: number }]
  update(data, key, onClick) {
    const container = document.getElementById(this.selector)
    const containerWidth = container.clientWidth
    // TODO: enforce a minimum height for label display
    const containerHeight = container.clientHeight

    this.svg.attr('width', containerWidth).attr('height', containerHeight)

    const plotWidth = containerWidth - this.margins.left - this.margins.right
    const plotHeight = containerHeight - this.margins.top - this.margins.bottom

    // update scales
    this.x.range([0, plotWidth]).domain([0, d3.max(data, d => d.value) * 1.1])
    this.y.range([0, plotHeight]).domain(data.map(d => d[key]))

    this.yAxisFn.scale(this.y)
    this.yAxis.call(this.yAxisFn)

    // join
    const bars = this.plot.selectAll('rect').data(data)
    const labels = this.plot.selectAll('text.label').data(data)
    const counts = this.plot.selectAll('text.count').data(data)

    // exit
    bars
      .exit()
      .transition(d3.transition())
      .remove()

    labels
      .exit()
      .transition(d3.transition())
      .remove()

    counts
      .exit()
      .transition(d3.transition())
      .remove()

    // enter / update
    bars
      .enter()
      .append('rect')
      .attr('fill', ' #1976d2')
      .on('mouseover', function() {
        d3.select(this).attr('opacity', 0.7)
      })
      .on('mouseleave', function() {
        d3.select(this).attr('opacity', 1)
      })
      .on('click', (_, d) => onClick(d))
      .style('cursor', 'pointer')
      .merge(bars)
      .transition(d3.transition())
      .attr('x', 0)
      .attr('y', d => this.y(d[key]))
      .attr('width', d => this.x(d.value))
      .attr('height', this.y.bandwidth())

    labels
      .enter()
      .append('text')
      .attr('class', 'label')
      .attr('font-size', '24px')
      .style('pointer-events', 'none')
      .merge(labels)
      .transition(d3.transition())
      .attr('fill', d =>
        d.name.length * 15 > this.x(d.value) ? 'black' : 'white',
      )
      .attr('x', d =>
        d.name.length * 15 > this.x(d.value) ? this.x(d.value) + 10 : 10,
      )
      .attr('y', d => this.y(d[key]) + 26)
      .text(d => d.name)

    counts
      .enter()
      .append('text')
      .attr('class', 'count')
      .attr('fill', 'white')
      .attr('font-size', '24px')
      .style('pointer-events', 'none')
      .merge(counts)
      .transition(d3.transition())
      .attr('x', d => this.x(d.value) - 50)
      .attr('y', d => this.y(d[key]) + 26)
      .text(d => d.value)
  }
}

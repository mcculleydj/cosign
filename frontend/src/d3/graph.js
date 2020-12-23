import * as d3 from 'd3'

function clamp(x, lo, hi) {
  return x < lo ? lo : x > hi ? hi : x
}

function drag(simulation, width, height) {
  function dragStarted(event) {
    if (!event.active) simulation.alphaTarget(0.3).restart()
    event.subject.fx = event.subject.x
    event.subject.fy = event.subject.y
  }

  function dragged(event) {
    event.subject.fx = clamp(event.x, 15, width - 15)
    event.subject.fy = clamp(event.y, 15, height - 15)
  }

  function dragEnded(event) {
    if (!event.active) simulation.alphaTarget(0)
    event.subject.fx = null
    event.subject.fy = null
  }

  return d3
    .drag()
    .on('start', dragStarted)
    .on('drag', dragged)
    .on('end', dragEnded)
}

const colors = {
  bill: 'gray',
  democrat: '#1565c0',
  republican: '#f44336',
  other: 'green',
}

export default class Graph {
  constructor(selector) {
    this.selector = selector
    this.svg = d3.select(`#${selector}`).append('svg')
    this.graphGroup = this.svg.append('g')
    this.nodeGroup = this.graphGroup.append('g').selectAll('.node')
    this.linkGroup = this.graphGroup.append('g').selectAll('.link')

    const legend = this.svg.append('g')

    Object.entries(colors).forEach(([type, color], i) => {
      legend
        .append('circle')
        .attr('class', 'legend')
        .attr('fill', color)
        .attr('cx', 30)
        .attr('cy', 30 + 30 * i)
        .attr('r', 10)

      legend
        .append('text')
        .attr('class', 'legend-text')
        .attr('font-size', '16px')
        .attr('x', 50)
        .attr('y', 35 + 30 * i)
        .text(type.charAt(0).toLocaleUpperCase() + type.slice(1))
    })

    this.hoverCircle = legend
      .append('circle')
      .attr('class', 'hover-legend')
      .attr('cx', 30)
      .attr('cy', 160)
      .attr('r', 0)

    this.hoverText = legend
      .append('text')
      .attr('class', 'hover-legend-text')
      .attr('font-size', '16px')
      .attr('x', 50)
      .attr('y', 165)
  }

  update(graph) {
    const container = document.getElementById(this.selector)
    // TODO: min / max?
    const containerWidth = container.clientWidth
    const containerHeight = container.clientHeight

    // adjust the size of the SVG
    this.svg.attr('width', containerWidth).attr('height', containerHeight)
    this.graphGroup
      .attr('width', containerWidth)
      .attr('height', containerHeight)

    this.svg.call(
      d3
        .zoom()
        .extent([
          [0, 0],
          [containerWidth, containerHeight],
        ])
        .scaleExtent([0.33, 3])
        .on('zoom', ({ transform }) => {
          this.graphGroup.attr('transform', transform)
        }),
    )

    const showName = d => {
      this.hoverCircle
        .transition(d3.transition())
        .attr('r', 15)
        .attr('fill', colors[d.type])
      this.hoverText
        .transition(d3.transition())
        .text(
          d.display.length > 50 ? d.display.slice(0, 50) + '...' : d.display,
        )
    }

    const hideName = () => {
      this.hoverCircle.transition(d3.transition()).attr('r', 0)
      this.hoverText.transition(d3.transition()).text('')
    }

    this.nodeGroup = this.nodeGroup.data(graph.nodes, function(d) {
      return d.display
    })

    this.nodeGroup.exit().remove()

    this.nodeGroup = this.nodeGroup
      .enter()
      .append('circle')
      .attr('r', 10)
      .attr('fill', d => colors[d.type])
      .attr('class', 'node')
      .attr('class', 'graph-el')
      .style('cursor', 'pointer')
      .on('mouseover', function(_, d) {
        d3.select(this)
          .transition(d3.transition())
          .attr('r', 15)
        showName(d)
      })
      .on('mouseleave', function() {
        d3.select(this)
          .transition(d3.transition())
          .attr('r', 10)
        hideName()
      })
      .merge(this.nodeGroup)

    this.linkGroup = this.linkGroup.data(graph.links)

    this.linkGroup.exit().remove()

    this.linkGroup = this.linkGroup
      .enter('line')
      .append('line')
      .attr('stroke', 'gray')
      .attr('stroke-width', 1)
      .attr('class', 'link')
      .attr('class', 'graph-el')
      .merge(this.linkGroup)

    // ensures links appear behind nodes even as new data is added
    this.svg.selectAll('.graph-el').sort((a, b) => (a.id && !b.id ? 1 : -1))

    const simulation = d3
      .forceSimulation()
      .alpha(2)
      .nodes(graph.nodes)
      .force('charge', d3.forceManyBody().strength(-1000))
      .force(
        'radial',
        d3
          .forceRadial(
            Math.min(containerWidth, containerHeight) / 4,
            containerWidth / 2,
            containerHeight / 2,
          )
          .strength(d => (d.type === 'bill' ? 0.5 : 0)),
      )
      .force(
        'link',
        d3.forceLink(graph.links).id(d => d.id),
      )
      .force('center', d3.forceCenter(containerWidth / 2, containerHeight / 2))
      .on('tick', () => {
        this.nodeGroup.attr('cx', d => d.x).attr('cy', d => d.y)

        this.linkGroup
          .attr('x1', d => d.source.x)
          .attr('y1', d => d.source.y)
          .attr('x2', d => d.target.x)
          .attr('y2', d => d.target.y)
      })

    this.nodeGroup.call(drag(simulation, containerWidth, containerHeight))
  }
}

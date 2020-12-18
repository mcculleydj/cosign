import * as d3 from 'd3'

// set the dimensions and margins of the graph
const margin = { top: 20, right: 20, bottom: 20, left: 20 }
const width = 800 - margin.left - margin.right
const height = 500 - margin.top - margin.bottom

function clamp(x, lo, hi) {
  return x < lo ? lo : x > hi ? hi : x
}

function drag(simulation) {
  function dragstarted(event) {
    if (!event.active) simulation.alphaTarget(0.3).restart()
    event.subject.fx = event.subject.x
    event.subject.fy = event.subject.y
  }

  function dragged(event) {
    event.subject.fx = clamp(event.x, 0, width)
    event.subject.fy = clamp(event.y, 0, height)
  }

  function dragended(event) {
    if (!event.active) simulation.alphaTarget(0)
    event.subject.fx = null
    event.subject.fy = null
  }

  return d3
    .drag()
    .on('start', dragstarted)
    .on('drag', dragged)
    .on('end', dragended)
}

export function drawGraph(selector, graph, callbacks) {
  // let instance = constructGraph(member, memberMap, params.threshold)

  console.log(Object.keys(callbacks))

  const simulation = d3
    .forceSimulation(graph.nodes)
    .force(
      'link',
      d3.forceLink(graph.links).id(d => d.id),
    )
    .force('charge', d3.forceManyBody())
    .force('center', d3.forceCenter(width / 2, height / 2))

  const svg = d3
    .select(selector)
    .append('svg')
    .attr('viewBox', [0, 0, width, height])

  const link = svg
    .append('g')
    .attr('stroke', '#999')
    .attr('stroke-opacity', 0.6)
    .selectAll('line')
    .data(graph.links)
    .join('line')
    .attr('stroke-width', d => Math.sqrt(d.value))

  const node = svg
    .append('g')
    .selectAll('circle')
    .data(graph.nodes)
    .join('circle')
    .attr('r', 5)
    .call(drag(simulation))
    .classed('bill-node', d => d.type === 'bill')
    .classed('dem-node', d => d.type === 'democrat')
    .classed('rep-node', d => d.type === 'republican')
    .classed('oth-node', d => d.type === 'other')

  simulation.on('tick', () => {
    link
      .attr('x1', d => d.source.x)
      .attr('y1', d => d.source.y)
      .attr('x2', d => d.target.x)
      .attr('y2', d => d.target.y)

    node.attr('cx', d => d.x).attr('cy', d => d.y)
  })

  function remove() {
    svg.remove()
  }

  return remove
}

import * as d3 from 'd3'

export function constructGraph(member, memberMap, threshold) {
  const memberIDs = Object.keys(member.counts)
    .map(id => +id)
    .map(id => ({ ...memberMap[id], count: member.counts[id] }))
  const nodes = memberIDs
    .filter(m => m.count >= threshold)
    .concat(memberMap[member.id])
  const nodeIndex = {}
  nodes.forEach((m, i) => {
    nodeIndex[m.id] = i
  })

  return {
    nodes,
    links: Object.entries(member.counts)
      .map(([id, count]) => {
        if (count >= threshold) {
          return { source: nodeIndex[member.id], target: nodeIndex[id] }
        }
      })
      .filter(l => l),
  }
}

// set the dimensions and margins of the graph
const margin = { top: 20, right: 20, bottom: 20, left: 20 }
const width = 800 - margin.left - margin.right
const height = 500 - margin.top - margin.bottom

function clamp(x, lo, hi) {
  return x < lo ? lo : x > hi ? hi : x
}

export function drawGraph(selector, member, memberMap, params) {
  let instance = constructGraph(member, memberMap, params.threshold)

  const svg = d3
    .select(selector)
    .append('svg')
    .attr('width', width + margin.left + margin.right)
    .attr('height', height + margin.top + margin.bottom)
    .append('g')
    .attr('transform', 'translate(' + margin.left + ',' + margin.top + ')')

  const link = svg
    .selectAll('.link')
    .data(instance.links)
    .join('line')
    .classed('link', true)

  const node = svg
    .selectAll('.node')
    .data(instance.nodes)
    .join('circle')
    .attr('r', 6)
    .classed('node', true)
    .classed('dem', d => d.parties.length === 1 && d.parties.includes('D'))
    .classed('rep', d => d.parties.length === 1 && d.parties.includes('R'))
    .classed('oth', d => d.parties.length > 1)
    .classed('fixed', d => d.fx !== undefined)

  const simulation = d3
    .forceSimulation()
    .nodes(instance.nodes)
    .force('charge', d3.forceManyBody().strength(params.strength))
    .force('center', d3.forceCenter(width / 2, height / 2))
    .force('link', d3.forceLink(instance.links))
    .on('tick', tick)

  if (params.sticky) {
    const drag = d3
      .drag()
      .on('start', dragstart)
      .on('drag', dragged)
    node.call(drag).on('click', click)
  } else {
    const drag = d3
      .drag()
      .on('drag', dragged)
      .on('end', (_, d) => {
        delete d.fx
        delete d.fy
        simulation.alpha(1).restart()
      })
    node.call(drag)
  }

  function tick() {
    link
      .attr('x1', d => d.source.x)
      .attr('y1', d => d.source.y)
      .attr('x2', d => d.target.x)
      .attr('y2', d => d.target.y)
    node.attr('cx', d => d.x).attr('cy', d => d.y)
  }

  function click(_, d) {
    delete d.fx
    delete d.fy
    d3.select(this).classed('fixed', false)
    simulation.alpha(1).restart()
  }

  function dragstart() {
    d3.select(this).classed('fixed', true)
  }

  function dragged(event, d) {
    d.fx = clamp(event.x, 0, width)
    d.fy = clamp(event.y, 0, height)
    simulation.alpha(1).restart()
  }

  function reset() {
    node.dispatch('click')
  }

  function remove() {
    svg.remove()
  }

  return { reset, remove }
}

import * as d3 from 'd3'

// set the dimensions and margins of the graph
const margin = { top: 20, right: 20, bottom: 20, left: 20 }
const width = 400 - margin.left - margin.right
const height = 220 - margin.top - margin.bottom

// set the ranges
const x = d3.scaleLinear().range([0, width])
const y = d3.scaleLinear().range([height, 0])

// take a start angle in deg, an end angle in deg, and a number
// and provide an evenly spaced array of radians between the start and end
function tween(start, end, number) {
  let interval = (end - start) / number
  const degrees = []
  for (let i = 1; i <= number; i++) {
    degrees.push(start + interval * i)
  }
  return degrees.map(d => (d * Math.PI) / 180)
}

function applyPartyLabel(data, numDems, numReps, numOther) {
  // TODO: only works because there is a single non-Dem or Rep member
  // this should become more sophisticated for the general case
  if (numOther) {
    data[220].label = 'o'
  }
  let d = 0
  let r = 440
  while (numDems > 0 || numReps > 0) {
    if (numDems > 0 && d === 220 && numOther > 0) d++
    if (numDems > 0) data[d].label = 'd'
    if (numReps > 0) data[r].label = 'r'
    d++
    r--
    numDems--
    numReps--
  }
}

export function drawChart(selector, numDems, numReps, numOther) {
  const svg = d3
    .select(selector)
    .append('svg')
    .attr('width', width + margin.left + margin.right)
    .attr('height', height + margin.top + margin.bottom)
    .append('g')
    .attr('transform', 'translate(' + margin.left + ',' + margin.top + ')')

  // trial and error manual construction of a plot that
  // resembles the seating arrangement in the House
  let data = []

  const radii = []
  for (let i = 5; i < 12; i++) {
    radii.push(i)
  }

  radii.forEach(r => {
    tween(0, 37, r + 2).forEach(rad => {
      data.push({ x: r * Math.cos(rad), y: r * Math.sin(rad) })
      data.push({ x: -r * Math.cos(rad), y: r * Math.sin(rad) })
    })
  })

  tween(10, 35, 10).forEach(rad => {
    data.push({ x: 12 * Math.cos(rad), y: 12 * Math.sin(rad) })
    data.push({ x: -12 * Math.cos(rad), y: 12 * Math.sin(rad) })
  })

  radii.forEach(r => {
    tween(40, 70, r + 2).forEach(rad => {
      data.push({ x: r * Math.cos(rad), y: r * Math.sin(rad) })
      data.push({ x: -r * Math.cos(rad), y: r * Math.sin(rad) })
    })
  })

  tween(42, 68, 10).forEach(rad => {
    data.push({ x: 12 * Math.cos(rad), y: 12 * Math.sin(rad) })
    data.push({ x: -12 * Math.cos(rad), y: 12 * Math.sin(rad) })
  })

  radii.forEach(r => {
    tween(75, 105, r).forEach(rad => {
      data.push({ x: r * Math.cos(rad), y: r * Math.sin(rad) })
      data.push({ x: -r * Math.cos(rad), y: r * Math.sin(rad) })
    })
  })

  tween(75, 102, 11).forEach(rad => {
    data.push({ x: 12 * Math.cos(rad), y: 12 * Math.sin(rad) })
  })

  // reduce to 441 and sort by x
  data = data
    .sort((d1, d2) => d1.y - d2.y)
    .slice(2, data.length)
    .sort((d1, d2) => d1.x - d2.x)

  applyPartyLabel(data, numDems, numReps, numOther)

  // scale
  x.domain(d3.extent(data, d => d.x))
  y.domain(d3.extent(data, d => d.y))

  // draw
  svg
    .selectAll('dot')
    .data(data)
    .enter()
    .append('circle')
    .attr('r', 2)
    .attr('cx', d => x(d.x))
    .attr('cy', d => y(d.y))
    .style('stroke', ({ label }) => {
      if (label === 'd') return '#1565C0'
      else if (label === 'r') return '#F44336'
      else if (label === 'o') 'green'
      else return 'lightgray'
    })
    .style('fill', ({ label }) => {
      if (label === 'd') return '#1565C0'
      else if (label === 'r') return '#F44336'
      else if (label === 'o') return 'green'
      else return 'lightgray'
    })
}

export function clearChart(selector) {
  d3.select(selector)
    .select('svg')
    .remove()
}

export function memberToString(member) {
  if (!member) return ''
  return `${member.name} [${member.parties.join(' | ')} - ${
    member.state
  } - ${member.districts.join(' | ')}]`
}

export function truncate(s, len) {
  if (!s) {
    return ''
  }

  if (s.length > len) {
    const tokens = s.split(/\s+/)

    let prefix = ''
    let i = 0
    while (prefix.length < 25) {
      prefix = `${prefix} ${tokens[i]}`
      i++
    }

    let suffix = ''
    let j = tokens.length - 1
    while (suffix.length < 25) {
      suffix = `${tokens[j]} ${suffix}`
      j--
    }

    return `${prefix} ... ${suffix}`
  }

  return s
}

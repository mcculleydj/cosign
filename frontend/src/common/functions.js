export function memberToString(member) {
  if (!member) return ''
  return `${member.name} [${member.parties.join(' | ')} - ${
    member.state
  } - ${member.districts.join(' | ')}]`
}

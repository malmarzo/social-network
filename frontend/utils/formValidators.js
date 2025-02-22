// Form validator functions

export function validateTextOnly(input) {
  const regex = /^[A-Za-z]+$/;
  return regex.test(input);
}

export function validateNickname(input) {
  const regex = /^[A-Za-z0-9_]+$/;
  return regex.test(input);
}

export function validateDOB(input) {
  const dob = new Date(input);
  const today = new Date();
  let age = today.getFullYear() - dob.getFullYear();
  const monthDiff = today.getMonth() - dob.getMonth();
  if (monthDiff < 0 || (monthDiff === 0 && today.getDate() < dob.getDate())) {
    age--;
  }
  return dob < today && age >= 15;
}

export function validatePassword(input) {
  const regex =
    /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[A-Za-z\d@$!%*?&]{8,}$/;
  return regex.test(input);
}

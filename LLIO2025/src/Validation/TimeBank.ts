import { createForm } from 'felte';
import * as yup from 'yup';
import type { TimeBankConfig } from '../Models';

const schema = yup.object().shape({
  startDate: yup
    .string()
    .required("La date de début est requise"),

  hoursPerWeek: yup

    .number()
    .typeError("Veuillez entrer un nombre")
    .min(0, "Les heures ne peuvent pas être négatives")
    .required("Le nombre d'heures est requis"),

  offset: yup
    .number()
    .typeError("Veuillez entrer un nombre")
    .min(0, "Le décalage ne peut pas être négatif")
    .nullable()
});

export const validateTimeBankForm = (
  handleSubmit: (values) => void,
  config: TimeBankConfig
) => {
  return createForm({
    initialValues: { ...config },

    validate: async (values) => {
      try {
        await schema.validate(values, { abortEarly: false });
        return {};
      } catch (err) {
        const errors = {};
        err.inner.forEach((value) => {
          errors[value.path] = value.message;
        });
        return errors;
      }
    },

    onSubmit: handleSubmit
  });
};
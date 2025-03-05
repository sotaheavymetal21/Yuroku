import React, { TextareaHTMLAttributes, forwardRef } from 'react';
import ErrorMessage from './ErrorMessage';

interface TextareaFieldProps extends TextareaHTMLAttributes<HTMLTextAreaElement> {
  label?: string;
  error?: string;
  helperText?: string;
  fullWidth?: boolean;
}

const TextareaField = forwardRef<HTMLTextAreaElement, TextareaFieldProps>(
  (
    {
      label,
      error,
      helperText,
      fullWidth = true,
      className = '',
      id,
      rows = 4,
      ...rest
    },
    ref
  ) => {
    const inputId = id || `textarea-${Math.random().toString(36).substr(2, 9)}`;
    const widthClass = fullWidth ? 'w-full' : '';
    const errorClass = error
      ? 'border-red-500 focus:ring-red-500 focus:border-red-500'
      : 'border-gray-300 focus:ring-primary-500 focus:border-primary-500';

    return (
      <div className={`${widthClass} ${className}`}>
        {label && (
          <label htmlFor={inputId} className="block text-sm font-medium text-gray-700 mb-1">
            {label}
          </label>
        )}
        <textarea
          ref={ref}
          id={inputId}
          rows={rows}
          className={`appearance-none block px-3 py-2 border rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-2 ${errorClass} ${widthClass}`}
          {...rest}
        />
        {error && <ErrorMessage message={error} />}
        {helperText && !error && (
          <p className="mt-1 text-sm text-gray-500">{helperText}</p>
        )}
      </div>
    );
  }
);

TextareaField.displayName = 'TextareaField';

export default TextareaField; 